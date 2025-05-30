# Central Difference Method

## Overview
The central difference method approximates derivatives using points on both sides of the target point. It's generally more accurate than forward or backward differences for smooth functions.

## Mathematical Foundation

### First Derivative
For the first derivative, we use:
```
f'(x) ≈ [f(x+h) - f(x-h)] / (2h)
```

### Second Derivative
For the second derivative:
```
f''(x) ≈ [f(x+h) - 2f(x) + f(x-h)] / h²
```

### Third Derivative
For the third derivative:
```
f'''(x) ≈ [f(x+2h) - 2f(x+h) + 2f(x-h) - f(x-2h)] / (2h³)
```

## Error Analysis

### Linear Error O(h)
- Simple approximation with basic accuracy
- Suitable for preliminary calculations

### Quadratic Error O(h²)
- Standard accuracy for most applications
- Good balance between precision and computational cost

### Cubic Error O(h³)
- Higher accuracy using more function evaluations
- Better for sensitive calculations

### Quartic Error O(h⁴)
- Highest accuracy available
- Uses extended stencils for maximum precision

## Advantages
- Most accurate for interior points
- Symmetric formulation reduces bias
- Good stability properties

## Disadvantages
- Requires function evaluation at x±h
- Cannot be used at boundaries without modification
