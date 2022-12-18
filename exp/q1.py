class Solution:
    # @param A : string
    # @param B : string
    # @return an integer
    def solve(self, A, B):
        totals = 0
        for i in range(int(A), int(B) + 1):
            digits = list(str(i))
            
            max_digit = max([int(x) for x in digits ])
            min_digit = min([int(x) for x in digits])
            
            avg = (max_digit + min_digit) / 2
            
            xor = 0
            for digit in digits:
                xor = xor ^ int(digit)
                
            if xor > avg:
                totals += 1
        return totals % (10**9 + 7)
            
