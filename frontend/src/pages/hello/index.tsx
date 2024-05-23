import type {MouseEvent} from 'react'
import {useState} from 'react'

interface UserInfo {
	id: string
	username: string
	password: string
	age: string
	gender: string
}

const query = (
	value: string,
	setUserInfo: (value: UserInfo | ((prevState: UserInfo) => UserInfo)) => void,
) => {
	console.log(
		'import.meta.env.VITE_APP_API_QUERY',
		import.meta.env.VITE_APP_API_QUERY,
	)
	fetch(`${import.meta.env.VITE_APP_API_QUERY}?q=${value}`, {
		method: 'GET',
		headers: {'Content-type': 'application/json'},
	})
		.then(async (res) => {
			const data = await res.json()
			console.log(data)
			setUserInfo(data)
		})
		.catch((err) => console.log(err))
}

/**
 * @returns HelloPage JSX.Element
 */
export default function Page() {
	const [userInfo, setUserInfo] = useState<UserInfo>({
		id: 'string',
		username: 'string',
		password: 'string',
		age: 'string',
		gender: 'string',
	})

	return (
		<div>
			<table>
				<thead>
					<tr>
						<th>id</th>
						<th>username</th>
						<th>password</th>
						<th>age</th>
						<th>gender</th>
					</tr>
				</thead>
				<tbody>
					<tr>
						<td>{userInfo.id}</td>
						<td>{userInfo.username}</td>
						<td>{userInfo.password}</td>
						<td>{userInfo.age}</td>
						<td>{userInfo.gender}</td>
					</tr>
				</tbody>
			</table>

			<label htmlFor="query">
				<input
					type="button"
					name="query"
					id="query"
					value="Query"
					onClick={(e: MouseEvent<HTMLInputElement>) =>
						query(e.currentTarget.value, setUserInfo)
					}
				/>
			</label>
		</div>
	)
}
