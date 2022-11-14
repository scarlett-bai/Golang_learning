const defaultAvatarUrl = '/resources/account.png'
Page({
  data: {
    showPath: true,
    values: [
      {
        name: 'john',
        id: 1,
      },
      {
        name: 'mary',
        id: 2,
      },
      {
        name: 'tom',
        id: 3,
      }
    ],
    avatarUrl: defaultAvatarUrl,
    onChooseAvatar(e: any){
      const {avatarUrl} = e.detail
      this.setData({
        avatarUrl,
      })
    }
  },
})
 