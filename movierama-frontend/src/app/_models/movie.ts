export interface Movie {
  id: number;
  user_id: number;
  title: string
  description: string
  posted_by: string
  likes: number
  hates: number
  user_liked: boolean
  user_hated: boolean
  is_same_user: boolean
  time_ago: string
}
