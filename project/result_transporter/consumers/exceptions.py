from typhoon.components.result_transporter.executions.base_consumer import BaseConsumer
import datetime

class Exceptions(BaseConsumer):

    def __init__(self, config, loop):

        super().__init__(config, loop)
        self.collections = {
            "processor_exceptions": self.get_exception_collection("processor"),
            "fetcher_exceptions": self.get_exception_collection("fetcher"),
            "scheduler_exceptions": self.get_exception_collection("scheduler"),
            "exceptions": self.get_exception_collection("result_transporter"),
            "exception_types": self.get_exception_collection("types")
        }

    async def send(self, task):
        exception_id = self.get_exception_id(task)
        component = self.get_exception_component(task.queue_name)
        exception_message = self.get_exception_message(task)
        find = {
            "exception_id": exception_id,
            "taskid": task.task_id
        }
        await self.collections[task.queue_name].update(find,
                                                       {"$set": {
                                                           "task": self.prepare_task(task.task),
                                                           "taskid": task.task_id
                                                       },
                                                           "$inc": {"count": 1}}, upsert=True)

        await self.collections["exception_types"].update({"exception_id": exception_id,
                                                          "component": component
                                                          }, {
                                                             "$set": {
                                                                 "exception_id": exception_id,
                                                                 "component": component,
                                                                 "message": exception_message
                                                             },
                                                             "$inc": {"count": 1}
                                                         }, upsert=True)