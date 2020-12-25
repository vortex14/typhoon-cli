from typhoon.run import ConfigurationComponent
from typhoon.components.donor.donor import Donor
from typhoon.components.donor.api.api import DonorApi
from typhoon.components.donor.state.state import StateManager


if __name__ == "__main__":
    component = ConfigurationComponent("donor", Donor, DonorApi, StateManager)
    component.run()