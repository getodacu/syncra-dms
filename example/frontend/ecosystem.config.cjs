module.exports = {
    apps: [
        {
            name: "Syncra Frontend",
	    cwd: "/home/webapp/apps/syncra/frontend",
            script: "/home/webapp/apps/syncra/frontend/build/index.js",
            node_args: "--env-file=.env",
	    env: {
                NODE_ENV: "production",
                PORT: 3102,
                HOST: "0.0.0.0"
            },
            instances: 1,
            exec_mode: "fork",
            max_memory_restart: "300M",
            time: true
        }
    ]
};
