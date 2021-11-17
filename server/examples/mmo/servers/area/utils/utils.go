package utils

import (    
    "math"    
)

func CalcDistXZ(x0, z0, x1, z1 float32) float32 {
    off1 := x0 - x1
    off2 := z0 - z1
    return float32(math.Sqrt(float64(off1*off1 + off2*off2)) )
}

//func NewSimpleValidator()