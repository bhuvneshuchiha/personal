class TreeNode(object):
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


class Solution(object):
    def isSameTree(self, p, q):
        """
        :type p: Optional[TreeNode]
        :type q: Optional[TreeNode]
        :rtype: bool
        """
        lst1 = []
        def inorder(node):
            if not node:
                return
            inorder(node.left)
            lst1.append(node.value)
            inorder(node.right)
            return lst1

        lst = []
        def inorder_1(node):
            if not node:
                return
            inorder_1(node.left)
            lst.append(node.value)
            inorder_1(node.right)
            return lst

        lst1 = inorder(p)
        lst = inorder_1(q)

        if lst1 == lst:
            return True
        else:
            return False
