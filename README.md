# bspwm mark register

sxhkdrc:
```
#
# bspwm mark register
#

# mark active window
super + shift + {u, i, o, p}
    bmr mark -set {u, i, o, p},$(bspc query -N -n)

# focus marked window
super + {u, i, o, p}
   bmr mark -focus {u, i, o, p}
```
