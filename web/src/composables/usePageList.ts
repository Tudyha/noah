import { useRequest } from "vue-hooks-plus";
import type { PageResponse } from "@/types";

type UsePageListOptions<Q extends Record<string, any>> = {
  defaultPage?: number;
  defaultPageSize?: number;
  defaultFilters?: Q;
  manual?: boolean;
};

export function usePageList<T, Q extends Record<string, any> = Record<string, any>>(
  service: (params: { page: number; limit: number } & Q) => Promise<PageResponse<T>>,
  options: UsePageListOptions<Q> = {},
) {
  const currentPage = ref(options.defaultPage ?? 1);
  const pageSize = ref(options.defaultPageSize ?? 10);
  const filters = ref<Q>((options.defaultFilters as Q) ?? ({} as Q));

  const { data, loading, run } = useRequest(() =>
    service({
      page: currentPage.value,
      limit: pageSize.value,
      ...filters.value,
    }),
    { manual: options.manual ?? false },
  );

  const search = (form: Q) => {
    filters.value = form;
    currentPage.value = 1;
    run();
  };

  const changePage = (page: number, size?: number) => {
    if (size && size > 0) {
      pageSize.value = size;
    }
    currentPage.value = page;
    run();
  };

  const changePageSize = (size: number) => {
    if (size <= 0) return;
    pageSize.value = size;
    currentPage.value = 1;
    run();
  };

  return {
    data,
    loading,
    run,
    currentPage,
    pageSize,
    filters,
    search,
    changePage,
    changePageSize,
  };
}
