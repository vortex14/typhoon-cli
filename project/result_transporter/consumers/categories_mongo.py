

from typhoon.components.result_transporter.executions.base_consumer import BaseConsumer



class CategoriesMongo(BaseConsumer):

    def __init__(self, config, loop):

        super().__init__(config, loop)
        self.collection_us = self.get_mongo_collection(collection_name="categories.us", client_name="main")
        self.collection_ca = self.get_mongo_collection(collection_name="categories.ca", client_name="main")
        self.collections = {
            'us': self.collection_us,
            'ca': self.collection_ca
        } 

    async def sync_send(self, task):
        collection = self.collections[task.result["type_marketplace"]]
        for cat in task.result["categories"]:
            await collection.update_many({"id": cat["id"]}, {"$set": cat}, upsert=True)
       
    async def send(self, task):
        collection = self.collections[task.result["type_marketplace"]]
        for cat in task.result["categories"]:
            await collection.update_many({"id": cat["id"]}, {"$set": cat}, upsert=True)
        
    async def send_many(self, bucket):
        bulk_products = self.collection_products.initialize_ordered_bulk_op()
        for task in bucket:   
            
            bulk_products.find({"product_id": task.result["product_id"]}).upsert().update({"$set": task.result, "$inc": {"count": 1}})
        await bulk_products.execute()
        
        
