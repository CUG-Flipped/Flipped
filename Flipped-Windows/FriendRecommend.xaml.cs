using MaterialDesignThemes.Wpf;
using System;
using System.Collections.Generic;
using System.Linq;
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

        private void CollectBtn_Click(object sender, RoutedEventArgs e)
        {
            if (collectBtn.Kind == PackIconKind.StarBorder)
                collectBtn.Kind = PackIconKind.StarRate;
            else
                collectBtn.Kind = PackIconKind.StarBorder;
        }
    }
}
