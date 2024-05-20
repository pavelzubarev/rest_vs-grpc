import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    vus: 100, // number of virtual users
    duration: '30s', // duration of the test
};

export default function () {
    const url = 'http://localhost:5002/api/polygon/calculateArea';
    const payload = JSON.stringify({
        points: [
            { x: 0, y: 0 },
            { x: 4, y: 0 },
            { x: 4, y: 3 }
        ]
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    let res = http.post(url, payload, params);

    check(res, {
        'is status 200': (r) => r.status === 200,
        'response time < 200ms': (r) => r.timings.duration < 200,
    });

    sleep(1);
}
