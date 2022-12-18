class Solution:
    # @param A : string
    # @param B : list of integers
    # @return a list of integers
    def solve(self, A, B):
        results = list()
        
        for query in B:
            left = (query - 1) - 1
            right = (query-1) + 1
            count = 1
            while left >= 0 and right < len(A):
                if A[left] == A[right]:
                    left -= 1
                    right += 1
                    count += 2
                    
                else:
                    break
                
            results.append(count)
        return results

