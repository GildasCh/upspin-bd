var app = new Vue({
    el: '#app',
    data: {
        resource: resource,
        pages: pages,
        theaterImage: resource,
        number: 0,
        event: {}
    },
    created: function () {
        window.addEventListener('keyup', this.key);
    },
    methods: {
        key: function(keyEvent) {
            if (keyEvent.key == "ArrowLeft") {
                this.previous();
            } else if (keyEvent.key == "ArrowRight") {
                this.next();
            }
        },
        next: function () {
            this.number = (this.number + 1) % this.pages;
            this.theaterImage = this.resource + '?page=' + this.number;
        },
        previous: function () {
            this.number = (this.number - 1 + this.pages) % this.pages;
            this.theaterImage = this.resource + '?page=' + this.number;
        },
    }
})
