from typhoon.components.result_transporter.executions.base_consumer import BaseConsumer
from hashlib import md5
from datetime import datetime

class MongoPrice(BaseConsumer):

    def __init__(self, config, loop):

        super().__init__(config, loop)       
        
        self.collection_prices = self.get_mongo_collection(collection_name="prices", client_name="main")
        
    async def sync_send(self, task):
        #print(task.result)
        
        #return await self.send_to_another_project(project_name="r1_prepare", component="processor", task=task, callback_name="from_source")
        
        #key = self.get_key(task)
        result = await self.collection_prices.update_many({"price_id": task.result["price_id"]}, {"$set": task.result}, upsert=True)
        
        #print(result.raw_result)

    async def send(self, task):
        
        #await self.send_to_another_project(project_name="r1_prepare", component="processor", task=task, callback_name="from_source")
        #key = self.get_key(task)
        result = await self.collection_prices.update_many({"price_id": task.result["price_id"]}, {"$set": task.result, "$setOnInsert" : {"created_at" : datetime.datetime.utcnow()}}, upsert=True)

    async def send_many(self, bucket):
        bulk_products = self.collection_prices.initialize_ordered_bulk_op()
        for task in bucket:
            
            bulk_products.find({"identifier": task.result["identifier"]}).upsert().update({"$set": task.result, "$inc": {"count": 1}})
        await bulk_products.execute()
        
        
	#def get_price_id (self, task):
        #task.result["price_id"] = md5((str(task.result["identifier"]) + str(task.result["updated_at"])).encode()).hexdigest()
        #task.result["price_id"] = md5((("amazon.com" + str(task.result["type"]) + str(task.result["identifier"]) + str(datetime.utcnow()))).encode()).hexdigest()