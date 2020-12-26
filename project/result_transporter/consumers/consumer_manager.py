from project.result_transporter.consumers import mongo_production
from typhoon.extensions.result_transporter import exceptions
from project.result_transporter.consumers import mongo_queue
from project.result_transporter.consumers import mongo_prices
from project.result_transporter.consumers import categories_mongo

class ConsumerManager:

    def __init__(self, config, loop):
        self.config = config
        self.loop = loop
        self.mongo_production = mongo_production.MongoProduction(self.config, self.loop)
        self.exceptions = exceptions.Exceptions(self.config, self.loop)
        self.mongo_queue = mongo_queue.MongoQueue(self.config, self.loop)
        self.mongo_prices = mongo_prices.MongoPrice(self.config, self.loop)
        self.categories_mongo = categories_mongo.CategoriesMongo(self.config, self.loop)
        self.loop.create_task(self.create_indexes())
        
        self.consumers = {
            "mongo_production": {
                "consumer_instance": self.mongo_production,
                "bucket_required": False,
                "bucket_limit": 8
            },
            "exceptions": {
                "consumer_instance": self.exceptions,
                "bucket_required": False,
                "bucket_limit": 8
            },
            "R1": {
                "consumer_instance": None,
                "bucket_required": True
            },
            "mongo_queue": {
                "consumer_instance": self.mongo_queue,
                "bucket_required": False,
                "bucket_limit": 8
            },
            "mongo_prices": {
                "consumer_instance": self.mongo_prices,
                "bucket_required": False,
                "bucket_limit": 8
            },
            "categories_mongo": {
                "consumer_instance": self.categories_mongo,
                "bucket_required": False,
                "bucket_limit": 8
            }
        }

    async def create_indexes(self):
        queue = self.config.cache_storage.get_mongo_collection(collection_name="queue", client_name="main")
        await self.config.cache_storage.create_index_collection(queue, "queue")
        await self.config.cache_storage.create_compound_index(queue, ["identifier", "type"])
        
        products = self.config.cache_storage.get_mongo_collection(collection_name="products", client_name="main")
        await self.config.cache_storage.create_index_collection(products, "identifier")
        await self.config.cache_storage.create_index_collection(products, "product_id")
        
        prices = self.config.cache_storage.get_mongo_collection(collection_name="prices", client_name="main")
        await self.config.cache_storage.create_index_collection(prices, "identifier")
        await self.config.cache_storage.create_index_collection(prices, "updated_at")
        await self.config.cache_storage.create_index_collection(prices, "product_id")