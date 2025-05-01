from types import CoroutineType


class TreeNode:
    def __init__(self, val, left = None, right = None):
        self.val = val
        self.left = left
        self.right = right

    def __str__(self):
        return str(self.val)

A = TreeNode(1)
B = TreeNode(2)
C = TreeNode(3)
D = TreeNode(4)
F = TreeNode(5)
G = TreeNode(6)

A.left = B
A.right = C
B.left = D
B.right = F
C.left = G

print(A)
print("XXXXXXXXX")


def pre_order(node):
    if not node:
        return

    print(node)
    pre_order(node.left)
    pre_order(node.right)

def post_order(node):
    if not node:
        return

    pre_order(node.left)
    pre_order(node.right)
    print(node)

def inorder(node):
    if not node:
        return

    pre_order(node.left)
    print(node)
    pre_order(node.right)

pre_order(A)
print("XXXXXXXXX")
post_order(A)
print("XXXXXXXXX")
inorder(A)

