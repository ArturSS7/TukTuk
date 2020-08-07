var app = new Vue({
    el: '#app',
    data: {
        requests : []
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
        methods: {
            say: function (message) {
                alert(message)
            }
        }
    }
})