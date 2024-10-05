workspace {

    model {
        customer = person "Customer" "" "person"
        admin = person "Admin User" "" "person"

        emailSystem = softwareSystem "Email System" "Mailgun" "external"
        calendarSystem = softwareSystem "Calendar System" "Calendly" "external"

        taskManagementSystem  = softwareSystem "Task Management System"{
            webContainer = container "User Web UI" "" "" "frontend"
            adminContainer = container "Admin Web UI" "" "" "frontend"
            dbContainer = container "Database" "PostgreSQL" "" "database"
            apiContainer = container "API" "Go" {
                authComp = component "Authentication"
                crudComp = component "CRUD"
            }
        }

        # Relationships between people and software systems
        customer -> webContainer "Manages tasks"
        admin -> adminContainer "Manages users"
        apiContainer -> emailSystem "Sends emails"
        apiContainer -> calendarSystem "Creates tasks in Calendar"

        # Relationships between containers
        webContainer -> apiContainer "Uses"
        adminContainer -> apiContainer "Uses"
        apiContainer -> dbContainer "Persists data"

        # Relationships to/from components
        crudComp -> dbContainer "Reads from and writes to"
        webContainer -> authComp "Authenticates using"
        adminContainer -> authComp "Authenticates using"
    }

    views {
        systemContext taskManagementSystem {
            include *
            autolayout
        }

        container taskManagementSystem {
            include *
            autolayout
        }

        component apiContainer {
            include *
            autolayout
        }

        # Dynamic diagram can be used to showcase a specific feature or process
        dynamic taskManagementSystem "LoginFlow" {
            webContainer -> apiContainer "Sends login request with username and password"
            apiContainer -> webContainer "Returns JWT token"
            webContainer -> customer "Persists JWT token in local storage"
            autolayout
        }

        styles {
            element "Software System" {
                background #1168bd
                color #ffffff
            }

            element "person" {
                shape Person
            }

            element "external" {
                background #eeeeee
                border dashed
                color #000000
            }

            element "frontend" {
                shape WebBrowser
            }

            element "database" {
                shape Cylinder
            }
        }
    }
}