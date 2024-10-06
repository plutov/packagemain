workspace "formulosity" "Surveys as Code" {
    model {
        user = person "Responder" "Survey User" "user"
        admin = person "Admin User" "Console Admin" "user"

        vcsSystem = softwareSystem "Github" "Survey configuration repositories" "external"

        # Internal system
        surveysSystem  = softwareSystem "Formulosity Software" "Surveys as Code Platform" {
            surveyUIContainer = container "Survey UI" "Public survey pages" "Next.js" "frontend"
            adminUIContainer = container "Admin Console UI" "UI to manage surveys" "Next.js" "frontend"
            dbContainer = container "Database" "Surveys/Responses storage" "SQLite"
            backendContainer = container "Backend Service" "REST API" "Golang" {
                parserComp = component "Parser" "Parse surveys configurations"
                adminAPIComp = component "Admin API" "Endpoints to manage surveys"
                userAPIComp = component "User API" "Endpoints to answer surveys"
            }
        }

        # Relationships between people and software systems
        user -> surveyUIContainer "Load surveys and answer them"
        admin -> adminUIContainer "Manages surveys"

        # Relationships between containers
        surveyUIContainer -> backendContainer "Uses"
        adminUIContainer -> backendContainer "Uses"
        backendContainer -> dbContainer "Persists data"

        # Relationships to/from components
        adminUIContainer -> adminAPIComp "Manages surveys"
        surveyUIContainer -> userAPIComp "Answers surveys"
        vcsSystem -> parserComp "Fetches surveys configurations"
    }

    views {
        systemContext surveysSystem {
            include *
            autolayout
        }

        container surveysSystem {
            include *
            autolayout
        }

        component backendContainer {
            include *
            autolayout
        }

        # Dynamic diagram can be used to showcase a specific feature or process
        dynamic surveysSystem "SurveysParser" {
            vcsSystem -> backendContainer  "Fetches surveys configurations"
            backendContainer -> dbContainer "Persists parsed surveys"
            backendContainer -> surveyUIContainer "Load parsed surveys"
            autolayout
        }
    }
}
