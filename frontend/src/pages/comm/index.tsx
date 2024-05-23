import {useState} from 'react'
import {Box, Link, List, ListItem, Typography} from '@mui/joy'

/**
 * @returns JSXElement
 */
export default function Page() {
	const [serverList] = useState([
		{
			label: 'gitlab',
			url: 'http://159.75.231.54:7080/',
		},
		{
			label: 'argocd',
			url: 'http://159.75.231.54:6080/',
		},
		{
			label: 'alertmanager',
			url: 'http://159.75.231.54:9093',
		},
		{
			label: 'kiali',
			url: 'http://159.75.231.54:20001',
		},
		{
			label: 'prometheus',
			url: 'http://159.75.231.54:9090',
		},
		{
			label: 'grafana',
			url: 'http://159.75.231.54:3030',
		},
		{
			label: 'minio operator',
			url: 'http://159.75.231.54:9099',
		},
		{
			label: 'minio user',
			url: 'http://159.75.231.54:9443',
		},
		{
			label: 'jaeger',
			url: 'http://159.75.231.54:31821',
		},
		{
			label: 'harbor',
			url: 'http://159.75.231.54:30003',
		},
		{
			label: 'bitnami kafka',
			url: 'http://159.75.231.54:31476',
		},
		{
			label: 'ELK',
			url: 'http://159.75.231.54:5601',
		},
		{
			label: 'chaos',
			url: 'http://159.75.231.54:',
		},
		{
			label: 'consul',
			url: 'http://159.75.231.54:',
		},
		{
			label: 'higress',
			url: 'http://159.75.231.54:',
		},
		{
			label: 'casdoor',
			url: 'http://159.75.231.54:',
		},
	])
	const [noUiList] = useState([
		{
			label: 'gitlab',
			url: 'http://159.75.231.54:7080/',
		},
		{
			label: 'argocd',
			url: 'http://159.75.231.54:6080/',
		},
		{
			label: 'alertmanager',
			url: 'http://159.75.231.54:9093',
		},
		{
			label: 'kiali',
			url: 'http://159.75.231.54:20001',
		},
		{
			label: 'prometheus',
			url: 'http://159.75.231.54:9090',
		},
		{
			label: 'grafana',
			url: 'http://159.75.231.54:3030',
		},
		{
			label: 'minio operator',
			url: 'http://159.75.231.54:9099',
		},
		{
			label: 'minio user',
			url: 'http://159.75.231.54:9443',
		},
		{
			label: 'jaeger',
			url: 'http://159.75.231.54:31821',
		},
		{
			label: 'harbor',
			url: 'http://159.75.231.54:30003',
		},
		{
			label: 'bitnami kafka',
			url: 'http://159.75.231.54:31476',
		},
		{
			label: 'ELK',
			url: 'http://159.75.231.54:5601',
		},
		{
			label: 'chaos',
			url: 'http://159.75.231.54:',
		},
		{
			label: 'consul',
			url: 'http://159.75.231.54:31080',
		},
		{
			label: 'higress',
			url: 'http://159.75.231.54:',
		},
		{
			label: 'casdoor',
			url: 'http://159.75.231.54:8000',
		},
	])
	return (
		<Box>
			<Box>
				<Typography
					sx={{
						fontSize: '2.1rem',
						textAlign: 'center',
						color: '#fff',
					}}
				>
					已部署的微服务, 有UI界面:
				</Typography>
				<List
					sx={{
						ml: '15vw',
						width: '70vw',
						display: 'grid',
						gridTemplateColumns: '1fr 1fr 1fr',
						alignItems: 'center',
					}}
				>
					{serverList.map((item) => (
						<Link
							sx={{
								width: '20vw',
								height: '10vw',
								fontSize: '2rem',
								color: 'green',
								textAlign: 'center',
							}}
							key={item.label}
							href={item.url}
						>
							<ListItem
								sx={{
									width: '20vw',
									height: '10vw',
									background: '#3df',
								}}
							>
								{item.label}
							</ListItem>
						</Link>
					))}
				</List>
			</Box>
			<Box>
				<Typography
					sx={{
						color: '#fff',
						fontSize: '2.1rem',
						textAlign: 'center',
					}}
				>
					无UI界面的微服务:
				</Typography>
				<List
					sx={{
						ml: '15vw',
						width: '70vw',
						display: 'grid',
						gridTemplateColumns: '1fr 1fr 1fr',
						alignItems: 'center',
					}}
				>
					{noUiList.map((item) => (
						<ListItem
							sx={{
								width: '20vw',
								height: '10vw',
								textAlign: 'center',
								background: '#3df',
							}}
							key={item.label}
						>
							<Link
								sx={{
									color: 'green',
								}}
								href={item.url}
							>
								{item.label}
							</Link>
						</ListItem>
					))}
				</List>
			</Box>
		</Box>
	)
}
