import Table from '@mui/joy/Table'

interface RowData {
	name: string
	class?: string
	describe?: string
	state?: string
	note?: string
}

/**
 *
 * @param props 组件属性对象
 * @param props.title 标题
 * @param props.rows 行数据
 * @returns ReactElement 表格
 */
export default function TableCaption({
	title,
	rows,
}: {
	title: string
	rows: RowData[]
}) {
	return (
		<Table>
			<caption>{title}</caption>
			<thead>
				<tr>
					<th style={{width: '20%'}}>名称</th>
					<th style={{width: '40%'}}>类别</th>
					<th style={{width: '40%'}}>介绍</th>
					<th>状态</th>
					<th style={{width: '20%'}}>备注</th>
				</tr>
			</thead>
			<tbody>
				{rows.map((row) => (
					<tr key={row.name}>
						<td>{row.name}</td>
						<td>{row.class}</td>
						<td>{row.describe}</td>
						<td>{row.state}</td>
						<td>{row.note}</td>
					</tr>
				))}
			</tbody>
		</Table>
	)
}
