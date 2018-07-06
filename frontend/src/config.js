const dev = {
    api: {
        url: "http://localhost:8081",
    }
};

const prod = {
    api: {
        url: "/api",
    }
};

const config = process.env.REACT_APP_STAGE === 'prod'
    ? prod
    : dev;

export default {
    ...config
};