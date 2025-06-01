# Backward Difference Method

## Overview
The backward difference method approximates derivatives using points behind the target point. It's useful when forward points are not available or reliable.

## Mathematical Foundation

### First Derivative
For the first derivative, we use:
```
f'(x) ≈ [f(x) - f(x-h)] / h
```

### Second Derivative
For the second derivative:
```
f''(x) ≈ [f(x) - 2f(x-h) + f(x-2h)] / h²
```

### Third Derivative
For the third derivative:
```
f'''(x) ≈ [f(x) - 3f(x-h) + 3f(x-2h) - f(x-3h)] / h³
```

## Error Analysis

### Linear Error O(h)
- Basic backward difference formula
- Truncation error is proportional to h

### Quadratic Error O(h²)
- Uses additional backward points
- Improved accuracy through higher-order terms

### Cubic Error O(h³)
- Higher-order backward differences
- Better precision with more function evaluations

### Quartic Error O(h⁴)
- Maximum accuracy backward scheme
- Uses extended backward stencil

## Advantages
- Only requires function values at x and before
- Suitable for final boundary conditions
- Good for functions with singularities at x+h

## Disadvantages
- Generally less accurate than central differences
- Can accumulate errors in sequential calculations
- Requires careful handling of step size
