import re

def is_valid_name(name):
    pattern = r"^[A-Za-z]+(?:[-' ][A-Za-z]+)*$"
    return bool(re.fullmatch(pattern, name))

print(is_valid_name("John Doe"))
print(is_valid_name("Mary-Jane Smith"))
print(is_valid_name("O'Connor"))
print(is_valid_name("123John"))
print(is_valid_name("John@Doe"))

