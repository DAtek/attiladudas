import { notificationCollection, NotificationItem } from '@/components/notification/notification'
import type { errorMappingKey } from '@/utils/errors'
import { reactive } from 'vue'

export type PostTokenData = {
  username: string
  password: string
}

export type PostGalleryData = {
  title: string
  slug: string
  description: string
  date: string
  active: boolean
}

export type PutGalleryData = PostGalleryData & {
  id: number
}

export type PostFilesData = {
  formData: FormData
  galleryId: number
}

export type DeleteFilesData = {
  galleryId: number
  fileIds: number[]
}

export type GetGalleriesData = {
  page: number
  page_size: number
}

export type PostGalleryResult = {
  id: number
}

export type File = {
  id: number
  filename: string
  path: string
  rank: number
}

export type Gallery = {
  id: number
  title: string
  slug: string
  description: string
  active: boolean
  date: string
  files: File[]
}

export type GetGalleriesResult = {
  galleries: Gallery[]
  total: number
}

const KEY_TOKEN = 'token'

export type FieldError = {
  location: string
  type: errorMappingKey
  context: Object | null
  message: string
}

type ApiError = {
  errors: FieldError[]
}

export type Result<T> = {
  error: ApiError | undefined
  result: T | undefined
}

class ApiClient {
  protected postTokenUrl: () => URL
  protected postGalleryUrl: () => URL
  protected deleteGalleryUrl: (id: number) => URL
  protected putGalleryUrl: (id: number) => URL
  protected postFilesUrl: (galleryId: number) => URL
  protected deleteFilesUrl: (galleryId: number) => URL
  protected patchFileRankUrl: (fileId: number, rank: number) => URL
  protected getGalleryUrl: (slug: string) => URL
  protected getGalleriesUrl: () => URL
  protected _token: string

  constructor(protected onError?: (e: unknown) => void) {
    this.postTokenUrl = () => new URL('/api/token/', baseUrl())
    this.getGalleriesUrl = () => new URL('/api/galleries/', baseUrl())
    this.postGalleryUrl = () => new URL('/api/gallery/', baseUrl())
    this.putGalleryUrl = (id) => new URL(`/api/gallery/${id}/`, baseUrl())
    this.deleteGalleryUrl = (id) => new URL(`/api/gallery/${id}/`, baseUrl())
    this.postFilesUrl = (galleryId) => new URL(`/api/gallery/${galleryId}/files/`, baseUrl())
    this.deleteFilesUrl = (galleryId) => new URL(`/api/gallery/${galleryId}/files/`, baseUrl())
    this.patchFileRankUrl = (fileId: number, rank: number) => new URL(`/api/file/${fileId}/rank/${rank}/`, baseUrl())
    this.getGalleryUrl = (slug: string) => new URL(`/api/gallery/${slug}/`, baseUrl())
    this._token = window.localStorage.getItem(KEY_TOKEN) || ''
  }

  public get token() {
    return this._token
  }

  logOut() {
    this._token = ""
    window.localStorage.removeItem(KEY_TOKEN)
  }

  async postToken(data: PostTokenData): Promise<Result<null>> {
    const response = await this.fetch(this.postTokenUrl(), {
      ...sharedRequestParams,
      method: 'POST',
      body: JSON.stringify(data)
    })

    if (response.status !== 201)
      return {
        error: {
          errors: [{ location: '', type: 'REQUIRED', context: null, message: '' }]
        },
        result: undefined
      }
    const responseData = JSON.parse(await response.text())
    this._token = responseData.token
    window.localStorage.setItem(KEY_TOKEN, this._token)
    return { error: undefined, result: undefined }
  }

  async getGalleries(data: GetGalleriesData): Promise<Result<GetGalleriesResult>> {
    const url = new URL(this.getGalleriesUrl())
    url.searchParams.append('page', `${data.page}`)
    url.searchParams.append('page_size', `${data.page_size}`)
    const headers = { ...sharedRequestParams.headers }

    if (this._token) {
      // @ts-ignore
      headers['Authorization'] = `Bearer ${this._token}`
    }

    const response = await this.fetch(url, {
      ...sharedRequestParams,
      method: 'GET',
      headers
    })

    if (response.status !== 200) {
      throw 'Not implemented'
    }

    const result = JSON.parse(await response.text()) as GetGalleriesResult
    return { result, error: undefined }
  }

  async getGallery(slug: string): Promise<Result<Gallery>> {
    const url = this.getGalleryUrl(slug)

    const response = await this.fetch(url, {
      ...sharedRequestParams,
      method: 'GET'
    })

    if (response.status !== 200) {
      throw 'Not implemented'
    }

    const result = JSON.parse(await response.text()) as Gallery
    return { result, error: undefined }
  }

  async postGallery(data: PostGalleryData): Promise<Result<PostGalleryResult>> {
    const response = await this.fetch(this.postGalleryUrl(), {
      ...sharedRequestParams,
      method: 'POST',
      body: JSON.stringify(data),
      headers: {
        ...sharedRequestParams.headers,
        Authorization: `Bearer ${this._token}`
      }
    })

    let result
    let error
    switch (response.status) {
      case 201:
        result = JSON.parse(await response.text()) as PostGalleryResult
        return { error: undefined, result }
      case 400:
        error = JSON.parse(await response.text()) as ApiError
        return { error, result: undefined }
      default:
        throw 'Not implemented'
    }
  }

  async putGallery(data: PutGalleryData): Promise<Result<null>> {
    const body = JSON.stringify({
      title: data.title,
      description: data.description,
      date: data.date,
      active: data.active,
      slug: data.slug
    })

    const response = await this.fetch(this.putGalleryUrl(data.id), {
      ...sharedRequestParams,
      method: 'PUT',
      headers: {
        ...sharedRequestParams.headers,
        Authorization: `Bearer ${this._token}`
      },
      body
    })

    let error
    switch (response.status) {
      case 200:
        return { error: undefined, result: null }
      case 400:
        error = JSON.parse(await response.text()) as ApiError
        return { error, result: undefined }
      default:
        throw 'Not implemented'
    }
  }

  async deleteGallery(id: number): Promise<Result<null>> {
    const response = await this.fetch(this.deleteGalleryUrl(id), {
      ...sharedRequestParams,
      method: 'DELETE',
      headers: {
        ...sharedRequestParams.headers,
        Authorization: `Bearer ${this._token}`
      }
    })

    if (response.status !== 200) {
      throw 'Not implemented'
    }

    return { error: undefined, result: null }
  }

  async postFiles(data: PostFilesData): Promise<Result<null>> {
    const response = await this.fetch(this.postFilesUrl(data.galleryId), {
      ...sharedRequestParams,
      method: 'POST',
      headers: {
        Authorization: `Bearer ${this._token}`
      },
      body: data.formData
    })

    let error
    switch (response.status) {
      case 201:
        return { error: undefined, result: null }
      case 400:
        error = JSON.parse(await response.text()) as ApiError
        return { error, result: undefined }
      default:
        throw `Not implemented error handling for status ${response.status}`
    }
  }

  async patchFileRank(fileId: number, rank: number): Promise<Result<null>> {
    const response = await this.fetch(this.patchFileRankUrl(fileId, rank), {
      ...sharedRequestParams,
      method: 'PATCH',
      headers: {
        Authorization: `Bearer ${this._token}`
      }
    })

    switch (response.status) {
      case 200:
        return { error: undefined, result: null }
      default:
        throw `Not implemented error handling for status ${response.status}`
    }
  }

  async deleteFiles(data: DeleteFilesData): Promise<Result<null>> {
    const response = await this.fetch(this.deleteFilesUrl(data.galleryId), {
      ...sharedRequestParams,
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${this._token}`
      },
      body: JSON.stringify({
        ids: data.fileIds
      })
    })

    if (response.status !== 200) {
      return { result: undefined, error: { errors: [] } }
    }

    return { error: undefined, result: null }
  }

  protected async fetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
    try {
      return await fetch(input, init)
    } catch (e) {
      if (this.onError) this.onError(e)
      throw e
    }
  }
}

function baseUrl(): string {
  return window.location.origin
}

function onError() {
  notificationCollection.addItem(
    new NotificationItem('DANGER', 'Something unexpected happened. Check your network connection.')
  )
}
export const apiClient = reactive(new ApiClient(onError))

const sharedRequestParams: Partial<RequestInit> = {
  mode: 'same-origin',
  cache: 'no-cache',
  headers: { 'Content-Type': 'application/json' }
}
