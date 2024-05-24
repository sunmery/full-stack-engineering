import type {ChangeEvent} from 'react'
import {useEffect} from 'react'
import {useState} from 'react'

interface UserInfo {
	data: {
		id: number
		username: string
		password: string
		age: string
		gender: string
	}
	code: number
	message: string
}

/**
 * @returns HelloPage JSX.Element
 */
export default function Page() {
	const [name, setName] = useState<string>('lixia')
	const [userInfo, setUserInfo] = useState<UserInfo>({
		data: {
			id: 0,
			username: 'string',
			password: 'string',
			age: 'string',
			gender: 'string',
		},
		code: 400,
		message: 'NULL',
	})

	const query = (
		value: string,
		setUserInfo: (
			value: UserInfo | ((prevState: UserInfo) => UserInfo),
		) => void,
	) => {
		console.log(value)
		fetch(`${import.meta.env.VITE_APP_API_QUERY}?name=${value}`, {
			method: 'GET',
			headers: {'Content-type': 'application/json'},
		})
			.then(async (res) => {
				const data = await res.json()
				console.log(data)
				setUserInfo(data)
			})
			.catch((err) => {
				console.log(err)
			})
	}

	useEffect(() => {
		fetch(`${import.meta.env.VITE_APP_API_QUERY}`, {
			method: 'PUT',
			headers: {'Content-type': 'application/json'},
			body: JSON.stringify({
				name,
			}),
		})
			.then(async (res) => {
				const data = await res.json()
				console.log(data)
				setUserInfo(data)
			})
			.catch((err) => {
				console.log(err)
			})
	}, [])

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
						<td>{userInfo.data.id}</td>
						<td>{userInfo.data.username}</td>
						<td>{userInfo.data.password}</td>
						<td>{userInfo.data.age}</td>
						<td>{userInfo.data.gender}</td>
					</tr>
				</tbody>
			</table>

			<label htmlFor="query">
				<input
					name="query"
					id="query"
					value={name}
					onChange={(e: ChangeEvent<HTMLInputElement>) =>
						setName(e.target.value)
					}
				/>
			</label>

			<button onClick={() => query(name, setUserInfo)}>Query</button>
		</div>
	)
}
