import React from 'react'
import ReactDOM from 'react-dom/client'
import {createBrowserRouter, RouterProvider} from 'react-router-dom'

import App from '@/App'
import HelloPage from '@/pages/hello'

const router = createBrowserRouter([
	{
		path: '/',
		element: <App />,
	},
	{
		path: '/hello',
		element: <HelloPage />,
	},
])

ReactDOM.createRoot(document.getElementById('root')!).render(
	<React.StrictMode>
		<RouterProvider router={router} />
	</React.StrictMode>,
)
