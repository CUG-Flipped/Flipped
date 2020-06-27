using MaterialDesignThemes.Wpf;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Windows;
using System.Windows.Controls;
using System.Windows.Data;
using System.Windows.Documents;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Imaging;
using System.Windows.Shapes;

namespace Flipped_Win10
{
    /// <summary>
    /// FriendRecommend.xaml 的交互逻辑
    /// </summary>
    public partial class FriendRecommend : Window
    {
        public FriendRecommend()
        {
            InitializeComponent();
        }

        private void Window_MouseLeftButtonDown(object sender, MouseButtonEventArgs e)
        {
            DragMove();
        }

        private void CloseButton_Click(object sender, RoutedEventArgs e)
        {
            this.Close();
        }

        private void Information_Button_Click(object sender, RoutedEventArgs e)
        {
            Drawer.IsRightDrawerOpen = true;
        }

        private void BackButton_Click(object sender, RoutedEventArgs e)
        {
            Drawer.IsRightDrawerOpen = false;
        }

        private void Refresh_Button_Click(object sender, RoutedEventArgs e)
        {
            SetRecommendFriend();
        }

        private void Dislike_Button_Click(object sender, RoutedEventArgs e)
        {
            string targetUser = UserNameChip.Content as string;
            if (LocalDataBaseOperator.IsFriends(NetWork.UserPWD.Item1, targetUser))
                NetWork.DeleteFriend(targetUser);
        }

        private void Collect_Button_Click(object sender, RoutedEventArgs e)
        {
            if (collectBtn.Kind == PackIconKind.StarBorder)
                collectBtn.Kind = PackIconKind.StarRate;
            else
                collectBtn.Kind = PackIconKind.StarBorder;
        }

        private void Chat_Button_Click(object sender, RoutedEventArgs e)
        {

        }

        private void Like_Button_Click(object sender, RoutedEventArgs e)
        {
            string targetUser = UserNameChip.Content as string;
            if (!LocalDataBaseOperator.IsFriends(NetWork.UserPWD.Item1, targetUser))
                NetWork.AddFriend(targetUser);
        }

        private void Window_Loaded(object sender, RoutedEventArgs e)
        {
            LocalDataBaseOperator.createDB("friendList.db");
            LocalDataBaseOperator.addTable(NetWork.UserPWD.Item1);
            SetRecommendFriend();
        }

        private void SetRecommendFriend()
        {
            var result = NetWork.GetRecommendFriend();
            Debug.WriteLine(result.Item1);
            var user = result.Item2;
            var bitmapImage = new BitmapImage();
            bitmapImage.BeginInit();
            bitmapImage.CacheOption = BitmapCacheOption.None;
            bitmapImage.StreamSource = new MemoryStream(Convert.FromBase64String(user.Photo));
            bitmapImage.EndInit();
            recommendFriendAvtar.Source = bitmapImage;
            UserNameChip.Content = user.UserName;
            EmailChip.Content = user.Email;
            NameChip.Content = user.RealName;
            ProfessionChip.Content = user.Profession;
            AgeChip.Content = user.Age;
            RegionChip.Content = user.Region;
            HobbyChip.Content = user.Hobby;
        }
    }
}
