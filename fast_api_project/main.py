from enum import Enum
from typing import Annotated

from fastapi import FastAPI, Query
from pydantic import BaseModel

app = FastAPI()


class ModelClass(str, Enum):
    randall = "randall"
    porqueMaria = "porqueMaria"
    ragnar = "ragnar"


class Item(BaseModel):
    name: str
    description: str
    price: float
    tax: float | None


# @app.get("/items/{foo}")
# async def root(foo: int):
#     return {"message": f"Reading the url params {foo}"}
#


@app.get("/items/")
async def get_items(q: Annotated[str, Query(min_length=3)]):
    return {
        "message": "Hey there homie",
        "response": q if q else "Sorry no query params sent",
    }


@app.post("/create/item/{item_id}")
async def create_item(item: Item, item_id: int):
    result = item.model_dump()
    result.update({"item_id": item_id})
    return {"message": "Thanks for the request", "response": result}


@app.get("/models/{model_name}")
async def getModel(model_name: ModelClass):
    if model_name.value == "randall":
        return {"message": f"Model name is {model_name.value}"}
    if model_name.value == "porqueMaria":
        return {"message": f"Model name is {model_name.value}"}
    if model_name is ModelClass.ragnar:
        return {"message": f"Model name is {model_name.value}"}


@app.get("/files/{file_path:path}")
async def get_path(file_path: str):
    return {"message": "Thanks for sharing the path", "path": file_path}


@app.get("/bitch/{item_id}")
async def read_items(item_id: int, owner_id: str | None, needy: bool = False):
    return {
        "message": "Thanks bitch",
        "item_id": item_id if item_id else 0,
        "owner_id": owner_id,
        "needy": needy,
    }
