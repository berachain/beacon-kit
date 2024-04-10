import { getClient } from '@lodestar/api';
import { config } from '@lodestar/config/default';

export async function createClient() {
    const beaconNodeUrl = process.env.BEACON_NODE_URL;
    const client = getClient(
        { baseUrl: beaconNodeUrl, timeoutMs: 60_000 },
        { config }
    );

    {
        let r = await client.beacon.getGenesis();
        if (!r.ok) {
            throw r.error;
        }

        client.beacon.genesisTime = r.response.data.genesisTime;
    }

    {
        let r = await client.config.getSpec();
        if (!r.ok) {
            throw r.error;
        }

        client.beacon.secsPerSlot = r.response.data.SECONDS_PER_SLOT;
    }

    client.slotToTS = (slot) => {
        return client.beacon.genesisTime + slot * client.beacon.secsPerSlot;
    };

    return client;
}
