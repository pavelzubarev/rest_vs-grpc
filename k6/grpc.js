import grpc from 'k6/net/grpc';
import { check, sleep } from 'k6';

const client = new grpc.Client();
client.load(['../proto'], 'polygon.proto');

export let options = {
    vus: 100, // number of virtual users
    duration: '30s', // duration of the test
};

export default function () {
    client.connect('localhost:5001', {
        plaintext: true,
    });

    const data = {
        points: [
            { x: 0, y: 0 },
            { x: 4, y: 0 },
            { x: 4, y: 3 }
        ],
    };

    const response = client.invoke('polygon.PolygonService/CalculateArea', data);

    check(response, {
        'status is OK': (r) => r && r.status === grpc.StatusOK,
    });

    client.close();
    sleep(1);
}
