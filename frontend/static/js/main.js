var app = new Vue({
    el: '#app',
    data: {
        requests : [],
        dns_req : [],
        resp : false,
        tcp_req :[],
        isHidden: true,
        port: '',
        message: '',
        start_success: '',
        StopIsHidden: true,
        stop_port: '',
        stop_success: '',
        TcpIsHidden: true,
        tcp_port: '',
        RunningIsHidden: true,
        running_tcp: [],
    },
    methods : {
        getData: function (proto) {
            axios.get('/api/' + proto + '?limit=10')
                .then(response => {
                    this.requests = response.data
                    console.log(response);
                })
                .catch(error => {
                    // handle error
                    console.log(error);
                })
        },
        getDNS: function () {
            axios.get('/api/dns/available')
                .then(response => {
                    this.requests = response.data
                    console.log(response);
                })
                .catch(error => {
                    // handle error
                    console.log(error);
                })
        },
        generateDNS: function (time) {
            axios.post('/api/dns/new', {
                delete_time: time
            })
                .then(response => {
                    this.dns_req = response.data
                    this.resp = true
                    console.log(response);
                })
                .catch(error => {
                    // handle error
                    console.log(error);
                })
        },
        StartTCP: function (port, message) {
            axios.post('/api/tcp/new', {
                port: port,
                message: message
            })
                .then(response => {
                    console.log(response)
                    this.getDNS()
                })
        },
        StopTCP: function (port) {
            axios.post('/api/tcp/shutdown', {
                port: port,
            })
                .then(response => {
                    console.log(response)
                    if (response.data["success"] === true) {
                        console.log("true")
                        this.stop_success = true
                        this.StopIsHidden = true
                    } else {
                        this.stop_success = false
                        this.StopIsHidden = true
                    }
                })
        },
        GetTCP: function (port) {
            axios.get('/api/tcp/data?port=' + port)
                .then(response => {
                    this.tcp_req = response.data
                    console.log(response);
                })
                .catch(error => {
                    console.log(error);
                })
        },
        GetRunningTCP: function () {
            axios.get('/api/tcp/running')
                .then(response => {
                    this.RunningIsHidden = false
                    this.running_tcp = response.data
                    console.log(response);
                })
                .catch(error => {
                    console.log(error);
                })
        },
        deleteDNS: function (domain) {
            axios.post('/api/dns/delete', {
                domain: domain
            })
        }
    },
})