import { api } from '@/shared/api/base';

export interface FeedbackRequest {
  email: string;
  name: string;
  message: string;
}

export const feedbackApi = {
  send: (data: FeedbackRequest) => api.post('/feedback', data),
};
