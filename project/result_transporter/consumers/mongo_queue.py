from typhoon.components.result_transporter.executions.base_consumer import BaseConsumer
from datetime import datetime 
from copy import deepcopy

class MongoQueue(BaseConsumer):

    def __init__(self, config, loop):

        super().__init__(config, loop)
        self.collection = self.get_mongo_collection(collection_name="queue", client_name="main")

    async def sync_send(self, task):
        print (task.result)
        bulk = self.collection.initialize_ordered_bulk_op()
        for _identifier in task.result["identifiers"]:
            document = deepcopy(_identifier)
            document["updated_at"] = datetime.utcnow()
            document["queue"] = task.result.get("queue", 1)
            bulk.find(_identifier).upsert().update({"$set": document, 
                                             "$inc": {"count": 1},
                                            "$setOnInsert": {
                                                "created_at": datetime.utcnow() 
                                            }})
        await bulk.execute()
        
        

    async def send(self, task):
        #print (task.result,1)
        bulk = self.collection.initialize_ordered_bulk_op()
        for _identifier in task.result["identifiers"]:
            document = deepcopy(_identifier)
            document["updated_at"] = datetime.utcnow()
            document["queue"] = task.result.get("queue", 1)
            #print(document)
            if document["queue"]:                
                bulk.find(_identifier).upsert().update({"$set": document, 
                                             "$inc": {"count": 1},
                                            "$setOnInsert": {
                                                "created_at": datetime.utcnow() 
                                            }})
            else:
                bulk.find(_identifier).update({"$set": document})
        
        await bulk.execute()

    async def send_many(self, bucket):
        bulk = self.collection.initialize_ordered_bulk_op()
        for task in bucket:
            
            bulk.find({"identifier": task.result["identifier"], "type": task.result["type"]}).upsert().update({"$set": task.result, "$inc": {"count": 1}})
        await bulk.execute()
