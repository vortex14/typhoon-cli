from typhoon.components.result_transporter.executions.base_consumer import BaseConsumer


class MongoProduction(BaseConsumer):

    def __init__(self, config, loop):

        super().__init__(config, loop)
        self.collection_products = self.get_mongo_collection(collection_name="products", client_name="main")

    async def sync_send(self, task):
        #print(task.result)
        #print ('here')
        #print (task)
        
        #return await self.send_to_another_project(project_name="r1_prepare", component="processor", task=task, callback_name="from_source")
        
        #key = self.get_key(task)
        
        result = await self.collection_products.update_many({"product_id": task.result["product_id"]}, {"$set": task.result}, upsert=True)
        #print(result.raw_result)

    async def send(self, task):
        #print (task.result)
        #print ('here')
        #await self.send_to_another_project(project_name="r1_prepare", component="processor", task=task, callback_name="from_source")
        #key = self.get_key(task)
        result = await self.collection_products.update_many({"product_id": task.result["product_id"]}, {"$set": task.result, "$inc": {"count": 1}}, upsert=True)
        
    async def send_many(self, bucket):
        bulk_products = self.collection_products.initialize_ordered_bulk_op()
        for task in bucket:   
            
            bulk_products.find({"product_id": task.result["product_id"]}).upsert().update({"$set": task.result, "$inc": {"count": 1}})
        await bulk_products.execute()
        
        
