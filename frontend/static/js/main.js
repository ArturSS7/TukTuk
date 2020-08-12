var app = new Vue({
    el: '#app',
    data: {
        requests : [],
        dns_req : [],
        resp : false
    },
    methods : {
        getData: function (proto) {
            axios.get('/api/'+proto+'?limit=10')
                .then( response => {
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
                .then( response => {
                    this.requests = response.data
                    console.log(response);
                })
                .catch(error => {
                    // handle error
                    console.log(error);
                })
        },
        generateDNS: function () {
            axios.get('/api/dns/new')
                .then( response => {
                    this.dns_req = response.data
                    this.resp = true
                    console.log(response);
                })
                .catch(error => {
                    // handle error
                    console.log(error);
                })
        }
    },
})