# cc-set-intersection

Confidential computing example - Set-Intersection

[@ref - research gate](https://www.researchgate.net/publication/228762579_Private_and_threshold_set-intersection)

```
Kissner L , Song D . Private and threshold set-intersection[J]. Advances in Cryptology – Crypto ’, 2004.
```

## problem

Alice has Set A (S<sub>A</sub>)

Bob has Set B (S<sub>B</sub>)

Alice and Bob want to compute the set-intersection of
S<sub>A</sub> and S<sub>B</sub>, but they do not want to
leak any information about their set.

This protocol is used for confidential computing.

Alice and Bob both run this protocol with their set's as input,
at the end of this protocol, they both have the intersection of
their sets.

In this protocol, `Set` is defined as a collection of big numbers.
For other type of elements, we can map them to big numbers.

## protocol

The protocol is based on section 6.2 of [Private and threshold set-intersection](https://www.researchgate.net/publication/228762579_Private_and_threshold_set-intersection), named `Intersection-Mal`.
