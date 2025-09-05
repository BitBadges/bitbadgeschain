# Acknowledgements

This project incorporates several components that were forked from the [Osmosis](https://github.com/osmosis-labs/osmosis). We acknowledge and thank the Osmosis team for their excellent work on these modules.

## Forked Components

### Third-Party Utilities

The following directories contain utilities and testing frameworks originally developed by the Osmosis team:

-   **`third_party/apptesting/`** - Testing utilities and helpers for blockchain application testing
-   **`third_party/simulation/`** - Simulation framework for testing blockchain state transitions
-   **`third_party/mocks/`** - Mock objects for testing pool and module interfaces
-   **`third_party/osmomath/`** - Mathematical utilities and decimal handling
-   **`third_party/osmoutils/`** - General utility functions and helpers

### Core Modules

The following modules were forked from Osmosis and adapted for the BitBadges chain:

-   **`x/gamm/`** - Generalized Automated Market Maker (GAMM) module for liquidity pools
-   **`x/poolmanager/`** - Pool management system for handling different types of liquidity pools

## Modifications

While these components maintain their core functionality from Osmosis, they have been modified to:

-   Integrate with the BitBadges chain architecture
-   Support BitBadges-specific token types and operations
-   Adapt to the BitBadges consensus and governance mechanisms
-   Remove Osmosis-specific dependencies and configurations

## License

These forked components maintain their original licenses from Osmosis (Apache 2.0). Please refer to the individual source files for specific licensing information.

---

_This acknowledgement is provided in good faith to recognize the contributions of the Osmosis project to the broader blockchain ecosystem._
