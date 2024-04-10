import { createHash } from 'node:crypto';

// port of https://github.com/ethereum/go-ethereum/blob/master/beacon/merkle/merkle.go
export function verifyProof(root, index, proof, value) {
    let buf = value;

    proof.forEach((p) => {
        const hasher = createHash('sha256');
        if (index % 2n == 0n) {
            hasher.update(buf);
            hasher.update(p);
        } else {
            hasher.update(p);
            hasher.update(buf);
        }
        buf = hasher.digest();
        console.log('-> ', toHex(buf));
        index >>= 1n;
        if (index == 0n) {
            throw new Error('branch has extra item');
        }
    });

    console.log('    ^^^ root');

    if (index != 1n) {
        throw new Error('branch is missing items');
    }

    if (toHex(root) != toHex(buf)) {
        throw new Error('proof is not valid');
    }

    console.log('proof ok!');
    console.log('<-');
}

export function toHex(t) {
    return '0x' + Buffer.from(t).toString('hex');
}

export function log2(n) {
    return Math.ceil(Math.log2(Number(n))) || 1;
}
