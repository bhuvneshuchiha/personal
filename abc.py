arr = "abc"
strlst = list(arr)
lst = []

for i in range(len(arr)-1):
    for j in range(i, len(arr)-1):
        temp = ""
        temp = strlst[i]
        strlst[i] = strlst[j]
        strlst[j] = temp
    lst.append("".join(strlst))

print(lst)
