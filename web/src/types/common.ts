// 搜索栏配置项的通用接口
export type SearchItem = {
  key: string; // 对应搜索参数的字段名
  label: string;
  type: 'input' | 'select' | 'date-range'; // 组件类型
  placeholder?: string;
  options?: { value: string | number; label: string }[]; // select组件使用
}

export type SearchProps = {
  items: SearchItem[]
}

// 表格列配置项的接口
interface Column {
  key: string; // 对应数据的字段名
  label: string;
  render?: (row: any) => VNode; // 自定义渲染函数，用于复杂内容或操作列
  isSortable?: boolean; // 是否可排序
}

// Table组件的主Props接口
export type TableProps =  {
  title?: string; // 表格标题
  data: any[]; // 表格数据
  columns: Column[]; // 列配置
  search?: SearchItem[]; // 搜索栏配置
  pagination: {
    currentPage: number;
    pageSize: number;
    total: number;
    onPageChange: (page: number) => void;
    onPageSizeChange: (size: number) => void;
  };
  isLoading?: boolean; // 加载状态
}