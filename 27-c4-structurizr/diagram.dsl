workspace "formulosity" "Surveys as Code" {
    model {
        # Actors
        user = person "User" "Survey User" "user"
        admin = person "Admin" "Console Admin" "user"

        # External systems
        githubSystem = softwareSystem "GitHub" "Survey configurations" "external"

        # Internal system
        surveysSystem = softwareSystem "Formulosity" "Surveys as code" {
            surveyUIContainer = container "Survey UI" "Survey web pages" "Next.js" "frontend"
            consoleUIContainer = container "Console UI" "Survey management page" "Next.js" "frontend"
            dbContainer = container "Database" "Survey storage" "SQLite" "database"
            apiContainer = container "API" "REST API" "Go" {
                parserComp = component "Parser" "Parse/validate surveys"
                userAPIComp = component "User API" "Public user API"
                adminAPIComp = component "Admin API" "Private admin API"
            }
        }

        # Relationships
        user -> surveyUIContainer "Uses"
        admin -> consoleUIContainer "Uses"

        # Containers
        apiContainer -> dbContainer "Stores surveys in"
        surveyUIContainer -> apiContainer "Calls"
        consoleUIContainer -> apiContainer "Calls"

        # Components
        surveyUIContainer -> userAPIComp "Calls"
        consoleUIContainer -> adminAPIComp "Calls"
        parserComp -> githubSystem "Fetches surveys config"
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

        component apiContainer {
            include *
            autolayout
        }

        dynamic surveysSystem "SurveyParser" {
            apiContainer -> githubSystem "Fetches the surveys"
            apiContainer -> dbContainer "Stores the surveys in"
            surveyUIContainer -> apiContainer "Displays available surveys"
            autolayout
        }
    }
}
