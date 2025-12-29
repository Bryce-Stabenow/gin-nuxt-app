export const useLists = () => {
  const config = useRuntimeConfig()
  const apiUrl = config.public.apiUrl

  interface ListItem {
    name: string
    quantity: number
    checked: boolean
    details?: string
    added_by: string
    added_at: string
  }

  interface List {
    id: string
    user_id: string
    name: string
    description?: string
    items: ListItem[]
    shared_with: string[]
    created_at: string
    updated_at: string
  }

  interface CreateListRequest {
    name: string
    description?: string
  }

  interface UpdateListRequest {
    name?: string
    description?: string
  }

  interface AddListItemRequest {
    name: string
    quantity?: number
    details?: string
  }

  interface UpdateListItemCheckedRequest {
    index: number
    checked: boolean
  }

  /**
   * Get headers with cookie forwarding for server-side requests
   */
  const getHeaders = (): Record<string, string> => {
    const headers: Record<string, string> = {}
    if (process.server) {
      const requestHeaders = useRequestHeaders(['cookie'])
      if (requestHeaders.cookie) {
        headers.cookie = requestHeaders.cookie
      }
    }
    return headers
  }

  /**
   * Create a new list
   */
  const createList = async (name: string, description?: string): Promise<List> => {
    const body: CreateListRequest = { name }
    if (description) {
      body.description = description
    }

    return await $fetch<List>(`${apiUrl}/lists`, {
      method: 'POST',
      credentials: 'include',
      headers: getHeaders(),
      body,
    })
  }

  /**
   * Get all lists for the authenticated user
   */
  const getLists = async (): Promise<List[]> => {
    return await $fetch<List[]>(`${apiUrl}/lists`, {
      method: 'GET',
      credentials: 'include',
      headers: getHeaders(),
    })
  }

  /**
   * Get a single list by ID
   */
  const getList = async (id: string): Promise<List> => {
    return await $fetch<List>(`${apiUrl}/lists/${id}`, {
      method: 'GET',
      credentials: 'include',
      headers: getHeaders(),
    })
  }

  /**
   * Update a list
   */
  const updateList = async (id: string, updates: UpdateListRequest): Promise<List> => {
    return await $fetch<List>(`${apiUrl}/lists/${id}`, {
      method: 'PUT',
      credentials: 'include',
      headers: getHeaders(),
      body: updates,
    })
  }

  /**
   * Add an item to a list
   */
  const addListItem = async (listId: string, item: AddListItemRequest): Promise<List> => {
    return await $fetch<List>(`${apiUrl}/lists/${listId}/items`, {
      method: 'POST',
      credentials: 'include',
      headers: getHeaders(),
      body: item,
    })
  }

  /**
   * Update an item's checked state
   */
  const updateListItemChecked = async (
    listId: string,
    itemIndex: number,
    checked: boolean
  ): Promise<List> => {
    return await $fetch<List>(`${apiUrl}/lists/${listId}/items/checked`, {
      method: 'PUT',
      credentials: 'include',
      headers: getHeaders(),
      body: {
        index: itemIndex,
        checked,
      } as UpdateListItemCheckedRequest,
    })
  }

  return {
    createList,
    getLists,
    getList,
    updateList,
    addListItem,
    updateListItemChecked,
  }
}

