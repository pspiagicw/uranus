class Solution:
    # @param A : list of integers
    # @param B : list of list of integers
    # @return a list of integers
    def solve(self, A, B):
        
        results = list()
        
        for query in B:
            l1, r1, l2, r2 = query[0], query[1], query[2], query[3]
            
            l1, r1, l2, r2 = l1-1, r1-1,l2-1,r2-1
            
            x1 = A[l1]
            for i in A[l1+1:r1+1]:
                x1 = x1 & i
            
            x2 = A[l2]
            for i in A[l2+1:r2+1]:
                x2 = x2 & i
                
            results.append(x1 ^ x2)
        
        return results
