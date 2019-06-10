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

## base knowledge

### Additively Homomorphic Cryptosystem

> `Paillier’s cryptosystem`

> `TODO:: 找 golang 的 Paillier 库, 或者自行实现(推荐)`

the encryption function with public key _pk_:

- E<sub>pk</sub>(·)

The cryptosystem
supports the following two operations,
which can be performed without knowledge of the
private key:

- E<sub>pk</sub>(a + b) := E<sub>pk</sub>(a) +<sub>h</sub> E<sub>pk</sub>(b)

  - Given the encryptions of a and b,
    E<sub>pk</sub>(a) and E<sub>pk</sub>(b),
    we can efficiently
    compute the encryption of a + b

  - 加同态

- E<sub>pk</sub>(c · a) := c ×<sub>h</sub> E<sub>pk</sub>(a)

  - Given a constant c and the encryption of a, E<sub>pk</sub>(a),
    we can efficiently compute the encryption of c · a,

  - 标量乘同态

We also require that the homomorphic public-key cryptosystem support secure (n, n)-
threshold decryption.

### Polynomials Over Rings

`TODO:: 需要一人确定是否用到此部分`

### Commitment

`TODO:: 需要具体方案`

### Hash Function

> `SHA`

> `TODO:: 需要确定使用哪种hash`

- h(·)

  - a hash function from
    {0, 1}<sup>\*</sup> to {0, 1}<sup>_l_</sup> ( _l_ = lg(1 / ε) ),
    where ε is a probability parameter chosen to be negligable)

### Zero-Knowledge Proofs

> `TODO:: 确定是否均为标准 ZKP 流程`

- POPK{C = E<sub>_pk_</sub>(x)}

  - V: 拥有数字 _x_ 的密文 C = E<sub>_pk_</sub>(x)
  - P: 证明知道秘密 x

- ZKPK{_f_ | p' = _f_ ∗<sub>_h_</sub> E<sub>_pk_</sub>(p)}

  - V: 拥有多项式 _p_ 的密文 p' = _f_ ∗<sub>_h_</sub> E<sub>_pk_</sub>(p)
  - P: 证明知道多项式 _p_

- ZKPK{_f_ | (p' = _f_ ∗<sub>_h_</sub> E<sub>_pk_</sub>(p)) ∧ (y = E<sub>_pk_</sub> (_f_))}

  - V: 拥有多项式 _p_ 的密文 p' = _f_ ∗<sub>_h_</sub> E<sub>_pk_</sub>(p)
  - V: 拥有乘法系数 _f_ 的密文 y = E<sub>_pk_</sub> (_f_))
  - P: 证明知道多项式 _p_ 和 乘法系数 _f_

## protocol

> `TODO:: 理清流程, 将上方的内容对号入座`

The protocol is based on section 6.2 of [Private and threshold set-intersection](https://www.researchgate.net/publication/228762579_Private_and_threshold_set-intersection), named `Intersection-Mal`.

### input
