dct = {"name": "bhuvnesh", "age": 20}

def unpack_greet(name, age):
    print(name, age)
unpack_greet(**dct)  # --> This will get unpacked once the dict reaches function
def pack_greet(**abc):
    print(abc["name"], abc["age"])

pack_greet(
    name="bhuvnesh", age=20
)  # --> This will get packed once the argument reaches function
