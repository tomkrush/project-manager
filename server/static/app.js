function application() {
    return {
        show: false,
        user: {},
        issues: [],
        totalIssues: 0,
        userLoaded: false,
        issuesLoaded: false,
        search: '',
        fetch() {
            this.refreshMyself();
            this.refreshIssues();
        },
        refreshMyself() {
            this.userLoaded = false;
            fetch('/myself')
                .then(response => response.json())
                .then(data => {
                    this.user = data;
                    this.userLoaded = true;
                })
                .catch(error => {
                    console.error('Error fetching user:', error);
                });
        },
        refreshIssues() {
            this.issuesLoaded = false;
            
            fetch('/issues?text='+this.search)
                .then(response => response.json())
                .then(data => {
                    this.issues = data.issues;
                    this.totalIssues = data.total;
                    this.issuesLoaded = true;
                })
                .catch(error => {
                    console.error('Error fetching issues:', error);
                });
        }
    };
}

document.querySelectorAll('[x-component]').forEach(component => {
    const componentName = `x-${component.getAttribute('x-component')}`
    class Component extends HTMLElement {
        connectedCallback() {
            this.append(component.content.cloneNode(true))
        }

        data() {
            const attributes = this.getAttributeNames()
            const data = {}
            attributes.forEach(attribute => {
                data[attribute] = this.getAttribute(attribute)
            })
            return data
        }
    }
    customElements.define(componentName, Component)
});