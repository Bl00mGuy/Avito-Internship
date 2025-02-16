import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '10s', target: 50 },
        { duration: '20s', target: 200 },
        { duration: '10s', target: 0 },
    ],
};

export default function () {
    let authRes = http.post('http://localhost:8080/api/auth', JSON.stringify({
        username: "testuser",
        password: "testpassword"
    }), { headers: { 'Content-Type': 'application/json' } });

    check(authRes, { 'status is 200': (r) => r.status === 200 });

    let token = authRes.json("token");

    let buyRes = http.get('http://localhost:8080/api/buy/t-shirt', {
        headers: { 'Authorization': `Bearer ${token}` }
    });

    check(buyRes, { 'status is 200': (r) => r.status === 200 });

    sleep(1);
}