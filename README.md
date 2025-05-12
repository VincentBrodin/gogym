# GoGym

GoGym is a full-stack web application designed to help users manage and optimize their fitness routines. Built with Vue.js for the frontend, Go with Echo for the backend, and Bun for dependency management, GoGym offers features such as a calorie calculator, workout creator, workout sharing, detailed workout stats, weight control, rest management between sets, and more.

## Features

- **Calorie Calculator**: Estimate your daily calorie needs.
- **Workout Creator**: Design and schedule custom workouts.
- **Workout Sharing**: Share routines with friends or the community.
- **Workout Stats & Weight Control**: Monitor your progress with detailed statistics.
- **Rest Management**: Manage rest intervals between sets.
- **And More**: Continuously evolving to offer more tools for your fitness journey.

## Tech Stack

- **Frontend**: [Vue.js](https://vuejs.org/)
- **Backend**: [Go](https://golang.org/) with [Echo](https://echo.labstack.com/)

## Getting Started

### Prerequisites

Ensure you have the following installed:

- [Bun](https://bun.sh/)
- [Go](https://golang.org/)
- [Git](https://git-scm.com/)
- [Godo](https://github.com/vincentbrodin/godo/)

### Installation

1. **Clone the Repository**

   ```bash
   git clone https://github.com/VincentBrodin/gogym.git
   cd gogym
   ```

2. **Set Up the Application**

   ```bash
   godo init
   ```

3. **Start the Development Server**

   ```bash
   godo dev_back
   godo dev_front
   ```

   Access the app at the given url from bun/vite.

4. **Build for Production**

   ```bash
   godo build
   ```
