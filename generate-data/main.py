from os import name
from typing import Dict, List
import pandas
import csv

from pymongo import MongoClient
from typing import (
    Dict,
)
from dataclasses import dataclass

@dataclass
class Product:
    _id: int
    name: str
    description: str
    sku: str
    price: float

class SimpleMongoClient:
    def __init__(self, collection: str, data_path: str):
        self._cluser = MongoClient("localhost:27017", username="admin", password="admin")
        self._db = self._cluser[collection]
        self._collection = self._db[collection]
        self._path = data_path

    def insert_row(self, data: Dict):
        self._collection.insert_one(data)

    def parser_data(self, data: List[str], id: int):
        return Product(
            _id=id,
            name=data[0],
            description=data[1],
            sku=data[2],
            price=data[3]
        ).__dict__

    def write_data(self):
        with open(self._path) as csv_file:
            csv_reader = csv.reader(csv_file, )
            _ = next(csv_reader)
            linecount = 0
            for line in csv_reader:
                linecount += 1
                self.insert_row(self.parser_data(data=line, id=linecount))

def main():
    client = SimpleMongoClient("products", "./data/products.csv")
    client.write_data()

if __name__ == "__main__":
    main()
