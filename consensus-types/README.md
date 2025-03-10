# consensus-types

The various fork specific types purely for marshalling and ssz can be found in their respective packages, e.g. `deneb`.
These types should only be converted to for the purposes of marshalling and ssz.

The normal control flow should use types such as those in `blocks`.