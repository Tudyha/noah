// 搜索栏配置项的通用接口
export type SearchItem = {
  key: string; // 对应搜索参数的字段名
  label: string;
  type: 'input' | 'select' | 'date-range' | 'date'; // 组件类型
  placeholder?: string;
  options?: { value: string | number; label: string }[]; // select组件使用
  // default?: string | number | Date | Date[]
}

export type SearchProps = {
  items: SearchItem[]
}

// 表格列配置项的接口
export type TableColumn = {
  key: string; // 对应数据的字段名
  label: string;
  render?: (column:TableColumn, row: any, index: number) => VNode; // 自定义渲染函数，用于复杂内容或操作列
  width?: string;
  sortable?: boolean;
  slot?: string
}

// Table组件的主Props接口
export type TableProps = {
  title?: string; // 表格标题
  data: any[]; // 表格数据
  columns: TableColumn[]; // 列配置
  searchItems?: SearchItem[]; // 搜索栏配置
  total?: number;
  currentPage?: number;
  pageSize?: number;
  isLoading?: boolean; // 加载状态
}